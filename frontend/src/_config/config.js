let url = null;

const getURL = async ( subroute ) =>{
    return `${process.env.REACT_APP_API_HOST}/${subroute}`;
};

export {
    getURL
}