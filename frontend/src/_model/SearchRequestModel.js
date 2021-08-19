export class SearchRequestModel {
    searchPosts;
    searchUsers;
    query;
    hashtags;

    constructor({searchPosts, searchUsers, query, hashtags}){
        this.searchPosts = searchPosts;
        this.searchUsers = searchUsers;
        this.query = query;
        this.hashtags = hashtags;
    }
}

export default SearchRequestModel;