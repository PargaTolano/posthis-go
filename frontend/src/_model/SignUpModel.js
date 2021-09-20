export class SignUpModel{

    userName;
    tag;
    email;
    password;

    constructor({username, tag, email, password}){
        this.username = username;
        this.tag = tag;
        this.email = email;
        this.password = password;
    }
}

export default SignUpModel;