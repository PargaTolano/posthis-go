import {  getURL  }     from '_config';
import { authHeader, requestWrapper }   from '_helpers';

import RepostViewModel  from '_model/RepostViewModel';


/**
 * @param   {Number} id
 */
const getReposts = async (id) =>  requestWrapper( async ()=> fetch( await getURL( `api/reposts/${id}` ) ) );

/**
 * @param {Number | String} userId
 * @param {Number | String} postId
 */
const createRepost = async ( userId, postId ) => {

    const headers = authHeader();

    const options = {
        method: 'POST',
        headers
    };

    return requestWrapper( async () =>  fetch( await getURL( `api/reposts-create/${userId}/${postId}` ), options ) );
};

/**
 * @param {Number | String} id
 */
 const deleteRepost = async ( id ) => {

    const headers = authHeader();

    const options = {
        method: 'DELETE',
        headers
    };

    return requestWrapper( async () =>  fetch( await getURL( `api/reposts-delete/${id}` ), options ) );
};

export{
    getReposts,
    createRepost,
    deleteRepost
}