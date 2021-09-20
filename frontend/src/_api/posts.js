import {  getURL  }                         from '_config';
import { authHeader, requestWrapper }       from '_helpers';

const getPosts = async () => requestWrapper( async ()=> fetch( await getURL( `api/posts` ) ) );

/**
 * @param {Number | String} id 
 * @returns 
 */
const getPost = async ( id ) => {

    const headers = authHeader();

    const options = { 
        headers
    };

    return requestWrapper( async ()=> fetch( await getURL( `api/post/${id}` ), options ) );
};

/**
 * @param {CPostModel} model
 */
const createPost = async ( model ) => {

    const headers = authHeader();

    let fd = new FormData();
    
    fd.append( 'content', model.content );
    for( let file of model.files ){
        fd.append('files', file );
    };

    const options = {
        method: 'POST',
        body: fd,
        headers
    };

    return requestWrapper( async () => fetch( await getURL( `api/posts-create` ), options ) );
};

/**
 * @param {Number | String} id 
 * @param {Object} model
 */
const updatePost = async ( id, model ) =>{

    const headers = {
        ...authHeader()
    };

    let fd = new FormData();

    fd.append('content', model.content);

    for( let id of model.deleted){
        fd.append('deleted', id);
    }

    for( let file of model.files){
        fd.append('files', file);   
    }

    const options = {
        method: 'PUT',
        body: fd,
        headers
    };
    return requestWrapper( async () => fetch( await getURL( `api/posts-update/${id}` ), options ) );
};

/**
 * @param   {Number | String} id
 */
const deletePost = async ( id ) =>{
    const options = {
        method: 'DELETE',
        headers: authHeader()
    };
    
    return requestWrapper( async () => fetch( await getURL( `api/posts-delete/${id}` ), options ) );
};

const getFeed = async ( offset, limit ) => {
    let headers= authHeader();

    let options ={
        headers
    };
    
    return requestWrapper( async () => fetch( await getURL( `api/posts-feed/${offset}/${limit}` ), options ) );
};

const getUserFeed = async ( id, offset, limit ) => {
    let headers= authHeader();

    let options ={
        headers
    };

    return requestWrapper( async () => fetch( await getURL( `api/posts-feed/${id}/${offset}/${limit}` ), options ) );
};

export{
    getPosts,
    getPost,
    createPost,
    updatePost,
    deletePost,
    getFeed,
    getUserFeed
}