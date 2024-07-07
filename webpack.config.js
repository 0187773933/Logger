const path = require( 'path' );

// npm root -g
// npm install -g webpack webpack-cli mini-css-extract-plugin css-loader crypto-browserify

const globalModulesPath = '/usr/local/lib/node_modules';

const MiniCssExtractPlugin = require( path.join( globalModulesPath, 'mini-css-extract-plugin' ) );
const cssLoaderPath = path.join( globalModulesPath , 'css-loader' );
const crypto_browserify = path.join( globalModulesPath , 'crypto-browserify' );

module.exports = {
    entry: [
        './v1/server/cdn/jquery.min.js',
        './v1/server/cdn/bootstrap.bundle.min.js',
        './v1/server/cdn/bcrypt.min.js',
        './v1/server/cdn/utils.js',
        './v1/server/cdn/bootstrap.min.css'
    ],
    output: {
        filename: 'bundle.js',
        path: path.resolve( __dirname, 'v1', 'server', 'cdn' )
    },
    module: {
        rules: [
            {
                test: /\.css$/,
                use: [
                    {
                        loader: MiniCssExtractPlugin.loader,
                        options: {
                            publicPath: ''
                        }
                    },
                    {
                        loader: cssLoaderPath
                    }
                ]
            }
        ]
    },
    plugins: [
        new MiniCssExtractPlugin( { filename: 'bundle.css' } )
    ],
    // resolve: {
    //     fallback: {
    //         "crypto": crypto_browserify
    //     }
    // },
    devtool: 'source-map',
    mode: 'production'
};
