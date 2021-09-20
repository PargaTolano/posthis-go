import {  getURL  }         from '_config';
import { authHeader, requestWrapper }       from '_helpers'
import { FollowViewModel }  from '_model';

/**
 * @param   {Number | string} id
 */
const getFollowers= async ( id ) => {

    const headers = authHeader();

    const options = {
        headers
    };

    return requestWrapper( async () => fetch( await getURL( `api/follows/${id}` ), options ));
}

/**
 * @param   {Number | string} id
 */
const getFollowing = async ( id ) => {

    const headers = authHeader();

    const options = {
        headers
    };
    
    return requestWrapper( async () => fetch( await getURL( `api/follows-following/${id}` ), options ));
};

/**
 * @param {Number | string} followedId
 */
const createFollow = async ( followedId ) => {

    const headers = authHeader();

    const options = {
        method: 'POST',
        headers
    };
    
    return requestWrapper( async () => fetch( await getURL( `api/follows-create/${followedId}` ), options ));
};

/**
 * @param {Number | string} id
 */
const deleteFollow = async ( id ) =>{

    const headers = authHeader();

    const options = {
        method: 'DELETE',
        headers
    };
    
    return requestWrapper( async () => fetch( await getURL( `api/follows-delete/${id}` ), options ));
};

//API ENTRYPOINT CORRECTION DONE

export{
    getFollowers,
    getFollowing,
    createFollow,
    deleteFollow
}