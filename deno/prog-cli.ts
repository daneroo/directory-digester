import { MultiBar, Presets } from "npm:cli-progress";
import { format as formatSize } from "@std/fmt/bytes";

async function processFile(filePath: string, multibar: MultiBar) {
  const fileInfo = await Deno.stat(filePath);
  const fileSize = fileInfo.size; // in bytes
  const totalFormattedSize = formatSize(fileSize);

  const rate = 100_000_000; // bytes per second
  const steps = Math.ceil(fileSize / rate);

  const fileName = filePath.split("/").pop() || filePath;
  const fileBar = multibar.create(fileSize, 0, {
    filename: fileName,
    totalFormattedSize,
  });

  for (let i = 0; i < steps; i++) {
    await new Promise((resolve) => setTimeout(resolve, 1000)); // Simulate processing
    const bytesProcessed = Math.min(rate, fileSize - fileBar.value);
    fileBar.increment(bytesProcessed);
  }

  multibar.remove(fileBar); // Remove the file progress bar after processing
}

async function processDirectory(directoryPath: string, multibar: MultiBar) {
  const entries = [];
  for await (const entry of Deno.readDir(directoryPath)) {
    entries.push(entry);
  }

  entries.sort((a, b) => a.name.localeCompare(b.name));

  const dirName = directoryPath.split("/").pop() || directoryPath;
  const totalFormattedSize = `${entries.length} files`;
  const dirBar = multibar.create(entries.length, 0, {
    filename: dirName,
    totalFormattedSize,
  });

  for (const entry of entries) {
    const fullPath = `${directoryPath}/${entry.name}`;
    if (entry.isDirectory) {
      // multibar.log(`Processing directory: ${entry.name}\n`);
      await processDirectory(fullPath, multibar);
    } else if (entry.isFile) {
      // multibar.log(`Processing file: ${entry.name}\n`);
      await processFile(fullPath, multibar);
    }
    dirBar.increment(); // Increment the directory progress bar for each processed child
  }

  multibar.remove(dirBar); // Optionally remove the directory progress bar after processing
}

const multibar = new MultiBar(
  {
    clearOnComplete: false,
    hideCursor: true,
    format:
      " {bar} {percentage}% | ETA: {eta}s | {totalFormattedSize} | {filename}",
  },
  Presets.shades_classic
);

const rootPath = Deno.args[0] || "testDirectories/rootDir01/";
await processDirectory(rootPath, multibar);

multibar.stop();
