import {  axios, getURL  }              from '_config';
import { authHeader, requestWrapper }   from '_helpers';

import { SearchRequestModel }           from '_model';

/**
 * @param {SearchRequestModel} model 
 */
const getSearch = async (
    query, 
    searchPosts, 
    searchUsers, 
    offsetPost, 
    limitPost, 
    offsetUser, 
    limitUser) => {
    const headers= authHeader();
    const options ={ headers };
    const url = new URL(getURL( `search/${offsetPost}/${limitPost}/${offsetUser}/${limitUser}`));
    url.searchParams.set( 'search-posts', searchPosts );
    url.searchParams.set( 'search-users', searchUsers );
    url.searchParams.set( 'query', query );
    
    return requestWrapper(()=>axios.get( url.href, options));
}

//TODO MANAGE DATAMODEL RETURNED WITH REQUEST WRAPPER

export{ getSearch }