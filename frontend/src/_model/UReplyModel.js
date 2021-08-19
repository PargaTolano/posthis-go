export class UReplyModel
{
    content;
    deleted;
    files  ;

    constructor({content, deleted, files}){
        this.content    = content;
        this.deleted    = deleted;
        this.files      = files;
    }
}

export default UReplyModel;