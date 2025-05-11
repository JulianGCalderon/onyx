import { defineCollection } from "astro:content";

import { glob } from "astro/loaders";

const note = defineCollection({
  loader: glob({ pattern: "**/*.md", base: "./content" }),
});

export const collections = { note };
