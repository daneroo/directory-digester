// TODO(daneroo): move external dependencies to deps.ts
import { parse } from "https://deno.land/std@0.180.0/flags/mod.ts";
import { basename, join } from "https://deno.land/std@0.180.0/path/mod.ts";
// fs is for walk
// import * as mod from "https://deno.land/std@0.180.0/fs/mod.ts";

// crypto.digest on a file: see https://examples.deno.land/hashing/ by https://www.linolevan.com/blog/tea_xyz
import {
  crypto,
  toHashString,
} from "https://deno.land/std@0.180.0/crypto/mod.ts";

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
    // Calculate the sha256 digest of the file
    const file = await Deno.open(node.path, { read: true });
    const readableStream = file.readable;
    const fileHashBuffer = await crypto.subtle.digest(
      "SHA-256",
      readableStream
    );
    const fileHash = toHashString(fileHashBuffer);

    node.info.sha256 = fileHash;
    logVerbose(`digestNode(${node.path}) = ${node.info.sha256} (leaf)`);
  } else {
    // Calculate the sha256 digest of the children
    const digester = new TextEncoder();
    // for (const child of node.children) {
    //   digester.encode(child.sha256);
    // }
    const childSha256s = node.children.map((child) =>
      digester.encode(child.info.sha256)
    );
    const arrayHashBuffer = await crypto.subtle.digest("SHA-256", childSha256s);
    const arrayHash = toHashString(arrayHashBuffer);
    node.info.sha256 = arrayHash;
    logVerbose(`digestNode(${node.path}) = ${node.info.sha256} (node)`);
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
  // const isDirIndicator = " "; // leaf or directory
  // if node.Info.Mode.IsDir() {
  // 	isDirIndicator = "/" //fmt.Sprintf("/ (%d)", len(node.Children))
  // }
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

// lets log to stderr - TODO turn this into loglevels?
function log(...data: any[]): void {
  console.error(`${new Date().toISOString()} -`, ...data);
}

let globalVerboseFlag = false;
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

  log(`directory-digester root: ${rootDirectory}`); // TODO(daneroo): add version,buildDate

  const rootNode = await buildTree(
    rootDirectory,
    await Deno.stat(rootDirectory)
  );

  logVerbose(
    `-- built tree: ${rootNode.info.name} (${rootNode.children.length})\n`
  );

  if (json) {
    showTreeAsJson(rootNode);
  } else {
    showAsIndented(rootNode, 0, 0);
  }
}

await main();
