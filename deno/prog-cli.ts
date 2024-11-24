import { MultiBar, Presets, Format } from "npm:cli-progress";
const { TimeFormat } = Format;
import { format as formatSize } from "@std/fmt/bytes";
import yoctoSpinner from "yocto-spinner";

interface DirectoryStats {
  path: string; // full path to this entry
  totalBytes: number; // total size including children
  totalEntries: number; // total count including children
  children: DirectoryStats[]; // ordered array of child stats
  parent?: DirectoryStats; // pointer to parent directory (undefined for root)
}

function incrementProgressAndStats(
  multibar: MultiBar,
  dirStats: DirectoryStats,
  bytes: number
) {
  // increment the progress bar for this directory
  for (const bar of multibar.bars) {
    bar.increment(bytes);
  }
  // increment the stats for this directory and all parents
  let current: DirectoryStats | undefined = dirStats;
  while (current !== undefined) {
    current.totalBytes += bytes;
    current = current.parent;
  }
}

async function processFile(fileStats: DirectoryStats, multibar: MultiBar) {
  const fileInfo = await Deno.stat(fileStats.path);
  const fileSize = fileInfo.size; // in bytes
  const totalFormattedSize = formatSize(fileSize);

  const rate = 1_000_000_000; // bytes per second
  const steps = Math.ceil(fileSize / rate);

  const fileName = fileStats.path.split("/").pop() || fileStats.path;
  const fileBar = multibar.create(fileSize, 0, {
    filename: fileName,
    totalFormattedSize,
  });

  if (steps === 1) {
    const waitTime = (fileSize / rate) * 1000;
    await new Promise((resolve) => setTimeout(resolve, waitTime));
    incrementProgressAndStats(multibar, fileStats, fileSize);
  } else {
    for (let i = 0; i < steps; i++) {
      await new Promise((resolve) => setTimeout(resolve, 1000));
      const bytesProcessed = Math.min(rate, fileSize - fileBar.value);
      incrementProgressAndStats(multibar, fileStats, bytesProcessed);
    }
  }

  multibar.remove(fileBar); // Remove the file progress bar after processing
}

async function processDirectory(dirStats: DirectoryStats, multibar: MultiBar) {
  const dirName = dirStats.path.split("/").pop() || dirStats.path;
  const totalFormattedSize = `${dirStats.totalEntries} files/dirs, ${formatSize(
    dirStats.totalBytes
  )}`;
  const dirBar = multibar.create(dirStats.totalBytes, 0, {
    filename: dirName,
    totalFormattedSize,
  });

  if (dirStats.children.length > 0) {
    for (const childStats of dirStats.children) {
      if (childStats.children.length > 0) {
        await processDirectory(childStats, multibar);
      } else {
        // It's a file
        await processFile(childStats, multibar);
        // progress bars are only incremented in processFile
      }
    }
  }
  multibar.remove(dirBar);
}

async function buildDirectoryStats(
  path: string,
  multibar: MultiBar
): Promise<DirectoryStats> {
  const entries: Deno.DirEntry[] = [];
  const dirStats: DirectoryStats = {
    path,
    totalBytes: 0, // these will be updated below
    totalEntries: 0, // these will be updated below
    children: [],
  };

  const stat = await Deno.stat(path);

  if (stat.isFile) {
    dirStats.totalBytes = stat.size;
    dirStats.totalEntries = 1;
  } else {
    // Read and sort directory entries
    for await (const entry of Deno.readDir(path)) {
      entries.push(entry);
    }
    entries.sort((a, b) => a.name.localeCompare(b.name));

    // Process each entry
    for (const entry of entries) {
      const childPath = `${path}/${entry.name}`;
      const childStats = await buildDirectoryStats(childPath, multibar);
      childStats.parent = dirStats;

      dirStats.children.push(childStats);
      dirStats.totalBytes += childStats.totalBytes;
      dirStats.totalEntries += childStats.totalEntries;
    }
    dirStats.totalEntries += 1; // Count this directory itself?
  }

  // multibar.log(
  //   `Discovered: ${path} (${dirStats.totalEntries} entries, ${dirStats.totalBytes} bytes)\n`
  // );
  return dirStats;
}

async function processPhases(rootPath: string, multibar: MultiBar) {
  // Phase 1: Discovery
  const spinner = yoctoSpinner({ text: "Phase 1: Discovery" }).start();
  // multibar.log("Phase 1: Discovery started\n");
  const start1 = Date.now();
  const rootDirStats = await buildDirectoryStats(rootPath, multibar);
  const elapsed1 = Date.now() - start1;
  // multibar.log(`Phase 1: Discovery completed in ${elapsed1}ms\n`);
  spinner.success(`Phase 1: Discovery completed in ${elapsed1}ms`);

  // Phase 2: Processing
  multibar.log("Phase 2: Digest started\n");
  const start2 = Date.now();
  await processDirectory(rootDirStats, multibar);
  const elapsed2 = Date.now() - start2;
  multibar.log(`Phase 2: Digest completed in ${elapsed2}ms\n`);
}

// deno-lint-ignore no-explicit-any
function paddedTimeFormat(t: any, options: any, roundToMultipleOf: any) {
  const formatted = TimeFormat(t, options, roundToMultipleOf);
  return formatted.padStart(7, " ");
}

const multibar = new MultiBar(
  {
    clearOnComplete: true,
    hideCursor: true,
    // fps: 5.0, // default is 10
    format:
      " {bar} {percentage}% | ETA: {eta_formatted} | {totalFormattedSize} | {filename}",
    //  {eta} is just in seconds, and not formatted.
    // " {bar} {percentage}% | ETA: {eta}s | {totalFormattedSize} | {filename}",
    // pad the time format
    // TimeFormat:= function formatTime(t, options, roundToMultipleOf){}
    formatTime: paddedTimeFormat,
  },
  Presets.shades_classic
);

const rootPath = Deno.args[0] || "testDirectories/rootDir01/";

await processPhases(rootPath, multibar);

// I don;t get the last log message if I don't wait; seems related to fps
await new Promise((resolve) => setTimeout(resolve, 100));
multibar.stop();
