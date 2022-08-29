import {  axios, getURL  }     from '_config';
import { authHeader, requestWrapper }   from '_helpers';

import RepostViewModel  from '_model/RepostViewModel';


/**
 * @param   {Number} id
 */
const getReposts = async (id) => requestWrapper(()=>axios.get(`reposts/${id}`));

/**
 * @param {Number | String} userId
 * @param {Number | String} postId
 */
const createRepost = async ( userId, postId ) => {
    const headers = authHeader();
    const options = { headers };
    return requestWrapper(()=>axios.post(`reposts-create/${userId}/${postId}`, options));
};

/**
 * @param {Number | String} id
 */
 const deleteRepost = async ( id ) => {
    const headers = authHeader();
    const options = { headers };
    return requestWrapper(()=>axios.delete(`reposts-delete/${id}`, options));
};

export{
    getReposts,
    createRepost,
    deleteRepost
}