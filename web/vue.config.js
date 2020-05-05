module.exports = {
  // serve on any path - single-page app with no history, so we should be fine
  publicPath: './',

  chainWebpack: config => {
    // make it easier to adjust output
    config.plugin('html')
      .tap(args => {
        args[0].minify = false;
        return args;
      });
  },
};
