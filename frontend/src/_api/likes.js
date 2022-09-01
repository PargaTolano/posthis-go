import {  axios, getURL  }     from '_config';
import { authHeader, requestWrapper } from '_helpers';

import LikeViewModel from '_model/LikeViewModel';

/**
 * @param   {Number} id
 */
const getLikes= async ( id ) => requestWrapper(()=>axios.get(`likes/${id}`));

/**
 * @param {Number | string} userId
 * @param {Number | string} postId
 * @returns
 */
const createLike = async ( userId, postId ) => {
    const headers = authHeader();
    const options = { headers };
    return requestWrapper(()=>axios.post(`likes-create/${userId}/${postId}`, {}, options));
};

/**
 * @param {Number | string} id
 */
 const deleteLike = async ( id ) => {
    const headers = authHeader();
    const options = { headers };

    return requestWrapper(()=>axios.delete(`likes-delete/${id}`, options));
};

export{
    getLikes,
    createLike,
    deleteLike
}