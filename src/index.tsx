import { copyFile, mkdir, readFile, writeFile, cp } from "fs/promises";
import { glob } from "glob";
import process from "process";
import rehypeStringify from "rehype-stringify";
import { renderToStaticMarkup } from "react-dom/server";

import remarkParse from "remark-parse";
import remarkRehype from "remark-rehype";
import { unified } from "unified";
import path from "path";
import rehypeReact from "rehype-react";
import { JSX } from "react/jsx-runtime";
import production from "react/jsx-runtime";

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

  const mdContent = await readFile(srcFile, "utf8");

  const htmlContent = await unified()
    .use(remarkParse)
    .use(remarkRehype, { allowDangerousHtml: true })
    .use(rehypeReact, production)
    .process(mdContent);

  const page = Template(htmlContent.result);

  const pageRaw = renderToStaticMarkup(page);

  const parent = path.dirname(dstFile);
  await mkdir(parent, { recursive: true });

  await writeFile(dstFile, pageRaw);
}

console.log("Copying styles!");
await cp("styles", "public/styles", { recursive: true });

function Template(content: JSX.Element) {
  return (
    <html>
      <head>
        <link rel="stylesheet" href="styles/main.css" />
      </head>
      <body>{content}</body>
    </html>
  );
}
