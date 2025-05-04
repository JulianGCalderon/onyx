import commonjs from '@rollup/plugin-commonjs';
import { nodeResolve } from '@rollup/plugin-node-resolve';

export default {
  input: [
    "mathjax.js"
  ],
  output: {
    dir: 'bundle',
    format: 'iife'
  },
  plugins: [commonjs(), nodeResolve()]
};

