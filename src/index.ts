import { copyFile, mkdir, readFile, writeFile, cp } from "fs/promises";
import { glob } from "glob";
import process from "process";
import rehypeStringify from "rehype-stringify";
import remarkParse from "remark-parse";
import remarkRehype from "remark-rehype";
import { unified } from "unified";
import path from "path";

console.log("CWD:", process.cwd());

const files = await glob("content/**", { nodir: true });

console.log("Generating content!");
for (const srcFile of files) {
  const relFile = path.relative("content", srcFile);
  let dstFile = path.join("public", relFile);

  console.log("- File: ", srcFile);

  if (path.extname(srcFile) != ".md") {
    await copyFile(srcFile, dstFile);
    continue;
  }

  dstFile = path.format({ ...path.parse(dstFile), base: "", ext: ".html" });

  const data = await readFile(srcFile, "utf8");

  const html = await unified()
    .use(remarkParse)
    .use(remarkRehype, { allowDangerousHtml: true })
    .use(rehypeStringify)
    .process(data);

  const parent = path.dirname(dstFile);
  await mkdir(parent, { recursive: true });

  await writeFile(dstFile, html.value);
}

console.log("Copying styles!");
await cp("styles", "public/styles", { recursive: true });
