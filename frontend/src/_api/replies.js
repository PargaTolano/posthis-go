import {  getURL  }     from '_config';
import { authHeader, requestWrapper }   from '_helpers';

import { CReplyModel, UReplyModel } from '_model';

/**
 * @param   {Number} id
 */
const getReplies = async ( id ) => {

    const headers = authHeader();

    const options = {
        headers
    };

    return requestWrapper( async ()=> fetch( await getURL( `api/replies/${id}` ), options ) );
};

/**
 * @param {CReplyModel} model
 */
const createReply = async ( model ) => {

    const headers = authHeader();

    let fd = new FormData();
    fd.append('content', model.content  );

    for( let f of model.files ){
        fd.append('files', f );
    }

    const options = {
        method: 'POST',
        body: fd,
        headers
    };

    return requestWrapper( async ()=> fetch( await getURL( `api/replies-create/${model.userID}/${model.postID}` ), options ) );
};

/**
 * @param {Number}      id
 * @param {UReplyModel} model
 */
 const updateReply = async ( id, model ) => {

    const headers = authHeader();

    let fd = new FormData();
    fd.append('content', model.content  );

    for( let did of model.deleted ){
        fd.append('deleted', did );
    }

    for( let f of model.files ){
        fd.append('files', f );
    }

    const options = {
        method: 'PUT',
        body: fd,
        headers
    };

    return requestWrapper( async ()=> fetch( await getURL( `api/replies-update/${id}` ), options ) );
};

/**
 * @param {Number} id
 */
 const deleteReply = async ( id ) => {

    const headers = authHeader();

    const options = {
        method: 'DELETE',
        headers
    };

    return requestWrapper( async ()=> fetch( await getURL( `api/replies-delete/${id}` ), options ) );
};

export{
    getReplies,
    createReply,
    updateReply,
    deleteReply
}