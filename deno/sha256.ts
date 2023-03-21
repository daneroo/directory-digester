// crypto.digest on a file: see https://examples.deno.land/hashing/ by https://www.linolevan.com/blog/tea_xyz
//  deno run --allow-read sha256.ts
const message = "The easiest, most secure JavaScript runtime.";

const messageBuffer = new TextEncoder().encode(message);

const hashBuffer = await crypto.subtle.digest("SHA-256", messageBuffer);

import {
  crypto,
  toHashString,
} from "https://deno.land/std@0.175.0/crypto/mod.ts";

// import { toHashString } from "https://deno.land/std@0.175.0/crypto/to_hash_string.ts";
const hash = toHashString(hashBuffer);
console.log({ hash });

const file = await Deno.open("example.txt", { read: true });

const readableStream = file.readable;

const fileHashBuffer = await crypto.subtle.digest("SHA-256", readableStream);

const fileHash = toHashString(fileHashBuffer);
console.log({ fileHash });
