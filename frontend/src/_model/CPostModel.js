export class CPostModel {
    userID;
    content;
    files;

    constructor({userID,content,files}){
        this.userID     = userID;
        this.content    = content;
        this.files      = files;
    }
}

export default CPostModel;