import { glob } from "glob";
import process from "process";

console.log("CWD:", process.cwd());

const files = await glob("content/**");
console.log("FILES:", files);
