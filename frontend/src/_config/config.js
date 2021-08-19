let url = null;

const getURL = async ( subroute ) =>{
    if( url === null){
        const res  = await fetch('/api.json');
        const { protocol, host, port} = await res.json();
    
        url = `${protocol}://${host}${port.length > 0 ? `:${port}`: port }`;
    }
    return `${url}/${subroute}`;
};

export {
    getURL
}