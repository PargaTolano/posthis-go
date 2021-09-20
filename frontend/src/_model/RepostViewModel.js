export class RepostViewModel{
    userID;
    postID;

    constructor({userID, postID}){
        this.userID = userID;
        this.postID = postID;
    }
}

export default RepostViewModel;