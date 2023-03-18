import { walkSync } from "https://deno.land/std/fs/mod.ts";
import { sha256 } from "https://deno.land/std/hash/mod.ts";

interface FileInfo {
  name: string;
  mod_time: string;
  mode: string;
  sha256?: string;
}

interface DirInfo {
  name: string;
  mod_time: string;
  mode: string;
  children?: FileInfo[];
  sha256?: string;
}

// Define the directory to walk recursively
// const root = "path/to/root/directory";
const root = "/Users/daniel/Downloads";

const dirStack: DirInfo[] = [];

for (const entry of walkSync(root, { includeDirs: true })) {
  if (entry.isDirectory) {
    // If the entry is a directory, add it to the stack and skip its children
    const dirInfo: DirInfo = {
      name: entry.path,
      mod_time: entry.mtime.toISOString(),
      mode: entry.mode.toString(8).slice(-3),
    };
    dirStack.push(dirInfo);
  } else {
    // If the entry is a file, add it to the current directory's list of children
    const fileInfo: FileInfo = {
      name: entry.path,
      mod_time: entry.mtime.toISOString(),
      mode: entry.mode.toString(8).slice(-3),
    };

    // Open the file
    const file = await Deno.open(entry.path);

    // Calculate the sha256 digest of the file
    const hasher = sha256.create();
    await Deno.copy(file, hasher);
    fileInfo.sha256 = hasher.hex();

    // Add the FileInfo to the current directory's list of children
    const dir = dirStack[dirStack.length - 1];
    if (!dir.children) {
      dir.children = [];
    }
    dir.children.push(fileInfo);
  }
}

// Calculate the sha256 digest of the FileInfo JSON structures for each directory's children
for (let i = dirStack.length - 1; i >= 0; i--) {
  const dir = dirStack[i];

  // Sort the list of children by name
  if (dir.children) {
    dir.children.sort((a, b) => a.name.localeCompare(b.name));
  }

  // Encode the list of children as JSON and calculate its sha256 digest
  const childrenJson = JSON.stringify(dir.children);
  dir.sha256 = sha256(childrenJson).hex();

  // Encode the DirInfo struct as JSON and print it
  const dirInfoJson = JSON.stringify(dir);
  console.log(dirInfoJson);
}
