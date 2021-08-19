export class CReplyModel
{
    userID ;
    postID ;
    content;
    files  ;

    constructor({userID, postID, content, files}){
        this.userID     = userID;
        this.postID     = postID;
        this.content    = content;
        this.files      = files;
    }
}

export default CReplyModel;