// crypto.digest on a file: see https://examples.deno.land/hashing/ by https://www.linolevan.com/blog/tea_xyz
//  deno run --allow-read sha256.ts

import { crypto } from "@std/crypto/crypto";
import { encodeHex } from "@std/encoding";

const message = "The easiest, most secure JavaScript runtime.";
const messageBuffer = new TextEncoder().encode(message);
const hashBuffer = await crypto.subtle.digest("SHA-256", messageBuffer);
const hash = encodeHex(new Uint8Array(hashBuffer));

console.log({ hash });

const file = await Deno.open("testDirectories/rootDir01/AA.txt", {
  read: true,
});

const readableStream = file.readable;

const fileHashBuffer = await crypto.subtle.digest("SHA-256", readableStream);
const fileHash = encodeHex(new Uint8Array(fileHashBuffer));
console.log({ fileHash });
