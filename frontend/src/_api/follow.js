import {  axios, getURL  }         from '_config';
import { authHeader, requestWrapper }       from '_helpers'
import { FollowViewModel }  from '_model';

/**
 * @param   {Number | string} id
 */
const getFollowers= async ( id ) => {
    const headers = authHeader();
    const options = { headers };

    return requestWrapper(()=>axios.get(`follows/${id}`, options));
}

/**
 * @param   {Number | string} id
 */
const getFollowing = async ( id ) => {
    const headers = authHeader();
    const options = { headers };
    
    return requestWrapper(()=>axios.get(`follows-following/${id}`, options));
};

/**
 * @param {Number | string} followedId
 */
const createFollow = async ( followedId ) => {
    const headers = authHeader();
    const options = { headers };
    
    return requestWrapper(()=>axios.post(`follows-create/${followedId}`, {}, options));
};

/**
 * @param {Number | string} id
 */
const deleteFollow = async ( id ) =>{
    const headers = authHeader();
    const options = { headers };
    
    return requestWrapper(()=>axios.delete(`follows-delete/${id}`, options));
};

//API ENTRYPOINT CORRECTION DONE

export{
    getFollowers,
    getFollowing,
    createFollow,
    deleteFollow
}