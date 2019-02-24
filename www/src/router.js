import createRouter from 'router5'
import browserPlugin from 'router5-plugin-browser';


export default function configureRouter() {
    const router = createRouter([], { allowNotFound: true})
        // Plugins
        router.usePlugin(
            browserPlugin({
               useHash: true
            })
        )
    return router
}

