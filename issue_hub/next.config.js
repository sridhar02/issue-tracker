require("dotenv").config();

module.exports = {
  env: {
    API_URL: process.env.API_URL
  },
  webpack: (config, { webpack }) => {
    config.plugins.push(
      new webpack.DefinePlugin({
        "process.env.API_URL": JSON.stringify(process.env.API_URL)
      })
    );

    return config;
  }
};
