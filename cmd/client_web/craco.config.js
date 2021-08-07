const CracoLessPlugin = require("craco-less");

module.exports = {
  plugins: [
    {
      plugin: CracoLessPlugin,
      options: {
        lessLoaderOptions: {
          lessOptions: {
            modifyVars: {
              // "@primary-color": "#1DA57A",
              "@input-height-base": "@height-base + 4px",
            },
            javascriptEnabled: true,
          },
        },
      },
    },
  ],
};
