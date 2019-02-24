import createRouter from 'router5'
// import browserPlugin from 'router5-plugin-browser';


export function configureRouter() {
    const router = createRouter([], { allowNotFound: true })
    // .usePlugin(
    //     browserPlugin()
    // )
    router.start()
    return router
}
