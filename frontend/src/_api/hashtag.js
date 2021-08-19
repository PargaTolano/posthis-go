import {  getURL  }     from '_config';

/**
 * @param   {Number} id
 */
const getHashtags= async ( ) => {
    let res = await fetch( await getURL( `api/hashtags/Get` ) );
    return res.json();
};

/**
 * @param   {Number} id
 */
 const getPostsWithHashtag = async ( text ) => {
    let res = await fetch( await getURL( `api/hashtags/GetPosts/${text}` ) );
    return res.json();
};

/**
 * @param {HashtagViewModel} model
 */
const createHashtag = async ( model ) => {

    const headers = {
        'Content-Type': 'application/json'
    }

    const options = {
        method: 'POST',
        body: JSON.stringify( model ),
        headers: headers
    };

    return fetch( await getURL( `api/hashtags/Create` ), options );
};

export{
    getHashtags,
    getPostsWithHashtag,
    createHashtag
}