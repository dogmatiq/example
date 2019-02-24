const path = require('path');
var HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
    mode: 'development',
    resolve: {
        extensions: ['.js', '.jsx']
    },
    // output: {
    //     path: path.resolve(__dirname, "dist"), // string
    // },
    module: {
        rules: [
            {
                test: /\.jsx?$/,
                loader: 'babel-loader'
            }
        ]
    },
    plugins: [new HtmlWebpackPlugin({
        template: './src/index.html'
    })],
    devServer: {
        overlay: true,
        compress: true,
        port: 9900
    },
    externals: {
        // global app config object
        config: JSON.stringify({
            apiUrl: 'http://localhost:9900'
        })
    }
}
