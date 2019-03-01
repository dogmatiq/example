var path = require('path');
var webpack = require('webpack');
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
    plugins: [
        new HtmlWebpackPlugin({
            template: './src/index.html'
        }),
        new webpack.DefinePlugin({
            'GPRCWEB_SERVER': JSON.stringify('http://localhost:8080')
        })
    ],
    devServer: {
        historyApiFallback: true,
        overlay: true,
        compress: true,
        port: 9900
    }
}
