import { axios }     from '_config';
import { authHeader, requestWrapper }   from '_helpers';

import { CReplyModel, UReplyModel } from '_model';

/**
 * @param   {Number} id
 */
const getReplies = async ( id ) => {
    const headers = authHeader();
    const options = { headers };
    return requestWrapper(()=>axios.get(`replies/${id}`, options));
};

/**
 * @param {CReplyModel} model
 */
const createReply = async ( model ) => {
    const headers = {
        ...authHeader(),
        'Content-Type': 'multipart/form-data'
    };

    const fd = new FormData();
    fd.append('content', model.content);
    model.files.forEach(x=>fd.append('files',x));

    const options = { headers };

    return requestWrapper(()=>axios.post(`replies-create/${model.userID}/${model.postID}`, fd, options));
};

/**
 * @param {Number}      id
 * @param {UReplyModel} model
 */
 const updateReply = async ( id, model ) => {
    const headers = {
        ...authHeader(),
        'Content-Type': 'multipart/form-data'
    };

    const fd = new FormData();
    fd.append('content', model.content);
    model.deleted.forEach(x=>fd.append('deleted', x));
    model.files.forEach(x=>fd.append('files', x));

    const options = { headers };

    return requestWrapper(()=>axios.put(`replies-update/${id}`, options));
};

/**
 * @param {Number} id
 */
 const deleteReply = async ( id ) => {
    const headers = authHeader();
    const options = { headers };

    return requestWrapper(()=>axios.delete(`replies-delete/${id}`, options));
};

export{
    getReplies,
    createReply,
    updateReply,
    deleteReply
}