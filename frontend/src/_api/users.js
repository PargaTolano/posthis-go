import {  getURL  }                 from '_config';
import { arrayToCSV }               from '_utils';
import { authHeader, requestWrapper }               from '_helpers';

import { SignUpModel, LogInModel, SearchRequestModel, UpdateUserViewModel} from '_model';

const getUsers = async () => {

    let headers= authHeader();

    let options ={
        headers
    };

    return requestWrapper( async ()=> fetch( await getURL( 'api/users-get' ), options ) );
}

const getUser = async ( id ) => {

    let headers= authHeader();

    let options ={
        headers
    };
    
    return requestWrapper( async ()=> fetch( await getURL( `api/user/${id}` ), options ) );
};

/**
 * @param {SignUpModel} model
 */
const createUser = async ( model ) => {

    const headers = {
        'Content-Type': 'application/json'
    };

    const options = {
        method: "POST",
        body: JSON.stringify( model ),
        headers
    };

    return requestWrapper( async ()=> fetch( await getURL( `api/users-create` ), options ) );
};


const validatePassword = async ( password ) => {

    const headers = authHeader();

    const options = {
        headers
    };

    return requestWrapper( async ()=> fetch( await getURL( `/api/validate-password/${password}`), options ) );
};

/**
 * @param {LogInModel} model
 */
const logIn = async ( model ) => {

    const headers = {
        'Content-Type': 'application/json'
    };

    const options = {
        method: 'POST',
        body: JSON.stringify( model ),
        headers
    };

    return requestWrapper( async ()=> fetch( await getURL( `api/login` ), options ) );
};

const logOut = async ( model ) => {

    const headers = authHeader();

    const options = {
        method: 'POST',
        headers
    };

    return requestWrapper( async ()=> fetch( await getURL( `api/logout` ), options ) );
};

/**
 * @param {Number} id 
 * @param {UpdateUserViewModel} model
 */
const updateUser = async ( id, model ) =>{

    const headers = authHeader();

    let body = new FormData();
    body.append('username'    , model.username    );
    body.append('tag'         , model.tag         );
    body.append('email'       , model.email       );
    body.append('profilePic'  , model.profilePic  );
    body.append('coverPic'    , model.coverPic    );

    const options = {
        method: 'PUT',
        body,
        headers
    };

    return requestWrapper( async ()=> fetch( await getURL( `api/users-update/${id}` ), options ) );
};

/**
 * @param   {Number} id
 */
const deleteUser = async ( id ) =>{
    const headers = authHeader();

    const options = {
        method: 'DELETE',
        headers
    };
    
    return requestWrapper( async ()=> fetch( await getURL( `api/users-delete/${id}` ), options ) );
};

export{
    getUsers,
    getUser,
    createUser,
    validatePassword,
    logIn,
    logOut,
    updateUser,
    deleteUser
}