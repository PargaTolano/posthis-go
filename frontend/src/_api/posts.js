import { axios } from '_config';
import { authHeader, requestWrapper } from '_helpers';

const getPosts = async () => requestWrapper(()=>axios.get(`posts`));

/**
 * @param {Number | String} id 
 * @returns 
 */
const getPost = async ( id ) => {
    const headers = authHeader();
    const options = { headers };

    return requestWrapper(()=>axios.get(`post/${id}`, options));
};

/**
 * @param {CPostModel} model
 */
const createPost = async ( model ) => {
    const headers = {
        ...authHeader(),
        'Content-Type': 'multipart/form-data'
    };

    const fd = new FormData();
    fd.append( 'content', model.content );
    model.files.forEach(x=>fd.append('files', x));

    const options = { headers };

    return requestWrapper(()=>axios.post(`posts-create`, fd, options));
};

/**
 * @param {Number | String} id 
 * @param {Object} model
 */
const updatePost = async ( id, model ) =>{
    
    const headers = {
        ...authHeader(),
        'Content-Type': 'multipart/form-data'
    };

    const fd = new FormData();
    fd.append('content', model.content);
    model.deleted.forEach(x=>fd.append('deleted',x));
    model.files.forEach(x=>fd.append('files', x));
    
    fd.forEach(x=>console.log(x));

    const options = { headers };

    return requestWrapper(()=>axios.put(`posts-update/${id}`, fd, options));
};

/**
 * @param   {Number | String} id
 */
const deletePost = async ( id ) =>{
    const headers = authHeader();
    const options = { headers };
    
    return requestWrapper(()=>axios.delete(`posts-delete/${id}`, options));
};

const getFeed = async ( offset, limit ) => {
    const headers= authHeader();
    const options = {headers};
    
    return requestWrapper(()=>axios.get(`posts-feed/${offset}/${limit}`, options));
};

const getUserFeed = async ( id, offset, limit ) => {
    const headers= authHeader();
    const options ={ headers };

    return requestWrapper(()=>axios.get(`posts-feed/${id}/${offset}/${limit}`, options));
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