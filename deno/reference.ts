// fs is for walk
// import * as mod from "https://deno.land/std@0.180.0/fs/mod.ts";
// crypto.digest on a file: see https://examples.deno.land/hashing/ by https://www.linolevan.com/blog/tea_xyz
import {
  crypto,
  toHashString,
} from 'https://deno.land/std@0.180.0/crypto/mod.ts';
// TODO(daneroo): move external dependencies to deps.ts
import { parse } from 'https://deno.land/std@0.180.0/flags/mod.ts';
import {
  basename,
  join,
} from 'https://deno.land/std@0.180.0/path/mod.ts';

// export VERSION=$(git describe --dirty --always)
// export COMMIT=$(git rev-parse --short HEAD)
// export BUILDDATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
// not sure how these would be packaged exposed
// when running from repo:url, (perhaps as a json file in the repo?)
// but they could be injected into a docker build as BUILD ARGS
// OR perhaps this? could work on github or npm,jsr??
// what if it was a jsonc file?
// const jsonURL = new URL('../package.json', import.meta.url);
// const jsonText = await Deno.readTextFile(jsonURL);
// const packageJson = JSON.parse(jsonText);

const buildInfo = {
  // version: "0.0.0-dev",
  version: Deno.env.get("VERSION") ?? "0.0.0-dev",
  commit: Deno.env.get("COMMIT") ?? "feedbac", // "c0ffee5"
  buildDate: Deno.env.get("BUILDDATE") ?? new Date().toISOString(),
};

interface DigestTreeNode {
  path: string;
  info: DigestInfo;
  children: DigestTreeNode[];
}

interface DigestInfo {
  name: string;
  size: number;
  mtime: Date;
  mode: number;
  sha256?: string;
}

function newDigestInfo(name: string, info: Deno.FileInfo): DigestInfo {
  const digestInfo: DigestInfo = {
    name: name,
    size: info.size,
    // mtime: info.mtime?.toISOString() ?? "",
    // panic if mtime is null?
    mtime: info.mtime || new Date(0),
    // panic if mode is null?
    mode: info.mode ?? -1,
  };
  return digestInfo;
}

function newLeaf(path: string, info: Deno.FileInfo): DigestTreeNode {
  const name = basename(path);
  const digestInfo: DigestInfo = newDigestInfo(name, info);
  return {
    path,
    info: digestInfo,
    children: [],
  };
}

// digestNode: calculates the digest of a node
// This can be invoked on a leaf node, or a directory node.
// On the directory it is assumed that the children have been previously digested
async function digestNode(node: DigestTreeNode): Promise<void> {
  // TODO(daneroo): should get from node.info.mode
  // const isFile = node.info.mode & 0o100000 === 0;
  const isFile = (await Deno.stat(node.path)).isFile;
  if (isFile) {
    const start = Date.now();

    // Calculate the sha256 digest of the file
    const file = await Deno.open(node.path, { read: true });
    const readableStream = file.readable;
    const fileHashBuffer = await crypto.subtle.digest(
      "SHA-256",
      readableStream
    );
    const fileHash = toHashString(fileHashBuffer);

    node.info.sha256 = fileHash;
    const elapsed = (Date.now() - start) / 1000;
    const sizeMB = node.info.size / (1024 * 1024);
    const rate = sizeMB / elapsed;

    logVerbose(
      `digestNode(${node.path}) = ${
        node.info.sha256
      } (leaf) - size: ${sizeMB.toFixed(2)}MB elapsed: ${elapsed.toFixed(
        2
      )}s rate: ${rate.toFixed(2)} MB/s`
    );
  } else {
    // Calculate the sha256 digest of the children
    const digester = new TextEncoder();
    const childSha256s = node.children.map((child) =>
      digester.encode(child.info.sha256)
    );
    const arrayHashBuffer = await crypto.subtle.digest("SHA-256", childSha256s);
    const arrayHash = toHashString(arrayHashBuffer);
    node.info.sha256 = arrayHash;
    // set size as sum of children's size (not the dir.stat.size)
    node.info.size = node.children.reduce(
      (sum, child) => sum + child.info.size,
      0
    );
    const sizeMB = node.info.size / (1024 * 1024);
    logVerbose(
      `digestNode(${node.path}) = ${
        node.info.sha256
      } (node) size: ${sizeMB.toFixed(2)}MB`
    );
  }
}

function ignoreName(name: string): boolean {
  const ignorePatterns = [".DS_Store", "@eaDir"];
  return ignorePatterns.some((pattern) => name.match(pattern) != null);
}

async function buildTree(
  parentPath: string,
  parentInfo: Deno.FileInfo
): Promise<DigestTreeNode> {
  logVerbose(`buildTree(${parentPath})`);
  const parentNode = newLeaf(parentPath, parentInfo);
  // need to sort the children, so we need to put then in an array
  // do not sort with localeCompare, as it is not
  const children: Deno.DirEntry[] = [];
  for await (const dirEntry of Deno.readDir(parentPath)) {
    children.push(dirEntry);
  }
  // children.sort((a, b) => a.name.localeCompare(b.name));
  children.sort((a, b) => (a.name > b.name ? 1 : a.name < b.name ? -1 : 0));

  for (const dirEntry of children) {
    const path = join(parentPath, dirEntry.name);
    const info = await Deno.stat(path);

    // ignore patterns
    const ignore = ignoreName(dirEntry.name);
    if (ignore) {
      logVerbose(`buildTree(${parentPath}) ignoring ${path}`);
      continue;
    }

    const node = info.isFile
      ? newLeaf(path, info)
      : await buildTree(path, info);

    if (info.isFile) {
      await digestNode(node);
    }
    parentNode.children.push(node);
  }

  // await setSizeOfParent(parentNode)
  await digestNode(parentNode);

  return parentNode;
}

function showAsIndented(
  node: DigestTreeNode,
  depth: number,
  maxLength: number
) {
  if (depth == 0 && maxLength == 0) {
    // maxLength = maxNameLength(node, 0)
    maxLength = 100;
  }
  const pad = " ".repeat(depth * 2);
  console.log(
    `${pad}${node.info.name} - ${node.info.size} bytes digest:${node.info.sha256}`
  );

  node.children.forEach((child) => {
    showAsIndented(child, depth + 1, maxLength);
  });
}

function convertTreeToListWithPath(node: DigestTreeNode, list: DigestInfo[]) {
  // Make a digestInfo with name replace by path
  const nameAsPathInfo = node.info;
  nameAsPathInfo.name = node.path;
  list.push(nameAsPathInfo);
  node.children.forEach((child) => {
    convertTreeToListWithPath(child, list);
  });
}

function showTreeAsJson(node: DigestTreeNode): void {
  const list: DigestInfo[] = [];
  convertTreeToListWithPath(node, list);
  // jsonBytes, err := json.MarshalIndent(list, "", "  ")
  console.log(JSON.stringify(list, null, 2));
}

// lets log to stderr - TODO turn this into log levels?
// deno-lint-ignore no-explicit-any
function log(...data: any[]): void {
  console.error(`${new Date().toISOString()} -`, ...data);
}

let globalVerboseFlag = false;
// deno-lint-ignore no-explicit-any
function logVerbose(...data: any[]): void {
  if (globalVerboseFlag) {
    log(...data);
  }
}

// write a function that builds a tree of files and directories
async function main() {
  // parse flags...
  const {
    verbose,
    json,
    // Define the directory to walk recursively
    // const root = "/Users/daniel/Downloads";
    _: [dirAsStrOrNum = "./go"],
  } = parse(Deno.args, {
    boolean: ["json", "verbose"],
    default: { verbose: false },
  });

  globalVerboseFlag = verbose;

  const rootDirectory = String(dirAsStrOrNum);

  // const buildInfo = JSON.parse(
  //   Deno.runSync({
  //     cmd: ["deno", "info", "--json"],
  //     stdout: "piped",
  //   }).stdout!
  // );

  // // Extract the relevant information
  // const gitCommit = buildInfo.gitCommit ?? "unknown-commit";
  // const buildTime =
  //   new Date(buildInfo.buildTime).toISOString() ?? "unknown-time";
  // const version = `1.0.0-${gitCommit}`;

  // These two lines are printed to stdout even if !verbose
  // TODO(daneroo) add a silent flag to suppress even these
  log(
    `directory-digester v${buildInfo.version} commit:${buildInfo.commit} build:${buildInfo.buildDate} deno:${Deno.version.deno}`
  ); // TODO(daneroo): add version,buildDate
  log(`directory-digester start root: ${rootDirectory}`); // TODO(daneroo): add version,buildDate
  const start = Date.now();

  const rootNode = await buildTree(
    rootDirectory,
    await Deno.stat(rootDirectory)
  );

  const elapsed = (Date.now() - start) / 1000; // convert to seconds
  const sizeMB = rootNode.info.size / (1024 * 1024);
  const rate = sizeMB / elapsed;

  // This line is printed to stdout even if !verbose
  // TODO(daneroo) add a silent flag to suppress even these
  log(
    `directory-digester done  root: ${rootNode.info.name} files: ${
      rootNode.children.length
    } - size: ${sizeMB.toFixed(2)}MB elapsed: ${elapsed.toFixed(
      2
    )}s rate: ${rate.toFixed(2)} MB/s`
  );
  if (json) {
    showTreeAsJson(rootNode);
  } else {
    showAsIndented(rootNode, 0, 0);
  }
}

await main();
