import { emailRegex, userNameRegex, passwordRegex, tagRegex } from '_utils';

const validateLogin = ({userName, password})=>{

    const validation = {
        userName:   true,
        password:   true,
        validated:  true
    };

    if ( !passwordRegex .test( password || '' ) ){
        validation.password = false;
        validation.validated = false;
    }
     

    if ( !( emailRegex.test( userName || '' ) || userNameRegex.test( userName || '' ) ) ){
        validation.userName = false;
        validation.validated = false;
    }
    
    return validation;
};

const validateSignup = ( {username, tag, email, password} ) =>{

    const validation = {
        username:   true,
        tag:        true,
        email:      true,
        password:   true,
        validated:  true,
    };

    if ( !userNameRegex .test( username || '' ) ){
        validation.username     = false;
        validation.validated    = false;
    }

    if ( !tagRegex .test( tag || '' ) ){
        validation.tag          = false;
        validation.validated    = false;
    }       

    if ( !emailRegex .test( email || '' ) ){
        validation.email        = false;
        validation.validated    = false;
    }

    if ( !passwordRegex .test( password || '' ) ){
        validation.password     = false;
        validation.validated    = false;
    }
    
    return validation;
};

const validateUpdateUser = ( {username, tag, email} ) =>{

    const validation = {
        username:   true,
        tag:        true,
        email:      true,
        validated:  true,
    };

    if ( !userNameRegex .test( username || '' ) ){
        validation.username     = false;
        validation.validated    = false;
    }

    if ( !tagRegex .test( tag || '' ) ){
        validation.tag          = false;
        validation.validated    = false;
    }       

    if ( !emailRegex .test( email || '' ) ){
        validation.email        = false;
        validation.validated    = false;
    }
    
    return validation;
};

const validateCreateAndUpdatePost = ( {content, mediaCount} )=>{

    const validation = {
        content:    true,
        mediaCount: true,
        validated:  true
    };

    if( content.length === 0 ){
        validation.content = false;
    }

    if( mediaCount === 0 ){
        validation.mediaCount = false;
    }

    validation.validated = validation.content || validation.mediaCount;

    return validation;
};

export{
    validateLogin,
    validateSignup,
    validateUpdateUser,
    validateCreateAndUpdatePost
}