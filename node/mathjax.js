import { mathjax } from "mathjax-full/js/mathjax.js";
import { TeX } from "mathjax-full/js/input/tex.js";
import { CHTML } from "mathjax-full/js/output/chtml.js";
import { liteAdaptor } from "mathjax-full/js/adaptors/liteAdaptor.js";
import { RegisterHTMLHandler } from "mathjax-full/js/handlers/html.js";
import { AssistiveMmlHandler } from "mathjax-full/js/a11y/assistive-mml.js";
import { AllPackages } from "mathjax-full/js/input/tex/AllPackages.js";

const FONT_URL =
  "https://cdn.jsdelivr.net/npm/mathjax@3/es5/output/chtml/fonts/woff-v2";

export default function() {
  const adaptor = liteAdaptor();
  const handler = RegisterHTMLHandler(adaptor);
  AssistiveMmlHandler(handler);
  const tex = new TeX({ packages: AllPackages });
  const chtml = new CHTML({ fontURL: FONT_URL });
  const html = mathjax.document("", { InputJax: tex, OutputJax: chtml });

  return {
    Css() {
      return adaptor.textContent(chtml.styleSheet(html));
    },

    Render(content, type) {
      const node = html.convert(content, {
        display: type == "display",
      });
      return adaptor.outerHTML(node);
    },
  };
}
