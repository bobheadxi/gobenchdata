module.exports = {
  // serve on any path - single-page app with no history, so we should be fine
  publicPath: './',

  css: {
    loaderOptions: {
      scss: {
        // @/ is an alias to src/
        prependData: `@import "~@/styles/variables.scss";`,
      },
    },
  },

  chainWebpack: config => {
    // make it easier to adjust output
    config.plugin('html')
      .tap(args => {
        args[0].minify = false;
        return args;
      });
  },
};
