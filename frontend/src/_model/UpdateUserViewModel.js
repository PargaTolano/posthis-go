export  class UpdateUserViewModel{

    userName     ;
    tag         ;
    email       ;
    profilePic  ;
    coverPic    ;

    constructor({username, tag, email, profilePic, coverPic}){
        this.username   = username   ;
        this.tag        = tag        ;
        this.email      = email      ;
        this.profilePic = profilePic ;
        this.coverPic   = coverPic   ;
    }

}

export default UpdateUserViewModel;