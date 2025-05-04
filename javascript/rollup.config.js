import commonjs from '@rollup/plugin-commonjs';
import { nodeResolve } from '@rollup/plugin-node-resolve';


export default {
  input: "main.js",
  output: {
    name: "mathjax",
    file: 'bundle.js',
    format: 'iife'
  },
  plugins: [commonjs(), nodeResolve()]
};

