import commonjs from '@rollup/plugin-commonjs';
import { nodeResolve } from '@rollup/plugin-node-resolve';
import { globSync } from 'glob';

export default {
  input: globSync("bin/*.js"),
  output: {
    dir: 'bundle',
    format: 'iife'
  },
  plugins: [commonjs(), nodeResolve()]
};

