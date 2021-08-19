export class LikeViewModel {
    postID;
    userID;

    constructor({postID, userID}){
        this.postID = postID;
        this.userID = userID;
    }
}

export default LikeViewModel;