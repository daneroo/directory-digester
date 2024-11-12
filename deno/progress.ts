import { format } from "@std/fmt/bytes";

// Simulate progress bars for each depth level using a simple object structure
const progressBars = new Map<string, { processed: number; total: number }>();

// deno-lint-ignore require-await
async function initializeProgress(name: string, total: number) {
  if (!progressBars.has(name)) {
    progressBars.set(name, { processed: 0, total });
    renderProgressBars();
  } else {
    const bar = progressBars.get(name);
    if (bar) {
      bar.total = total; // Update total if it changes dynamically
    }
  }
}

function updateProgress(name: string, steps: number = 1) {
  const bar = progressBars.get(name);
  if (bar) {
    bar.processed += steps;
    renderProgressBars();
  }
}

function renderProgressBars() {
  console.clear();
  progressBars.forEach((bar, name) => {
    const total = bar.total > 0 ? bar.total : 1; // Prevent division by zero
    const progressPercentage = Math.min(bar.processed / total, 1); // Ensure max 100%
    const filledLength = Math.min(Math.floor(progressPercentage * 20), 20); // Max 20
    const barString = "█".repeat(filledLength) + "-".repeat(20 - filledLength);
    console.log(
      `${name}: [${barString}] ${(progressPercentage * 100).toFixed(
        2
      )}% (${format(bar.processed)}/${format(bar.total)})`
    );

    // Remove completed progress bars
    if (progressPercentage >= 1) {
      progressBars.delete(name);
    }
  });
}

async function processFile(filePath: string, depth: number) {
  // Get the file size
  const fileInfo = await Deno.stat(filePath);
  const fileSize = fileInfo.size; // in bytes

  // Define the processing rate (e.g., 100 bytes per second)
  const rate = 2_000_000_000; // bytes per second

  // Calculate the number of steps based on the file size and rate
  const steps = Math.ceil(fileSize / rate);

  const fileName = filePath.split("/").pop() || filePath;
  await initializeProgress(fileName, fileSize);

  for (let i = 0; i < steps; i++) {
    // Simulate processing time for each step
    await new Promise((resolve) => setTimeout(resolve, 1000)); // 1 second per step
    // Update progress by the rate or remaining bytes
    const bytesProcessed = Math.min(rate, fileSize - i * rate);
    updateProgress(fileName, bytesProcessed);
  }
}

async function processDirectory(directoryPath: string, depth: number) {
  const entries = [];
  for await (const entry of Deno.readDir(directoryPath)) {
    entries.push(entry);
  }

  // Sort entries by name
  entries.sort((a, b) => a.name.localeCompare(b.name));

  // Initialize progress bar for the current directory
  const dirName = directoryPath.split("/").pop() || directoryPath;
  await initializeProgress(dirName, entries.length);

  for (const entry of entries) {
    const fullPath = `${directoryPath}/${entry.name}`;
    if (entry.isDirectory) {
      await processDirectory(fullPath, depth + 1);
    } else if (entry.isFile) {
      await processFile(fullPath, depth + 1);
    }
    updateProgress(dirName, 1);
  }
}

const rootPath = Deno.args[0] || "testDirectories/rootDir01/";
await processDirectory(rootPath, 0);
