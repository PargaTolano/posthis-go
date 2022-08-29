import {  getURL  } from '_config';
import { authHeader, requestWrapper } from '_helpers';
import { axios } from '_config';

import { 
    SignUpModel, 
    LogInModel, 
    SearchRequestModel, 
    UpdateUserViewModel
} from '_model';

const getUsers = async () => {
    const headers= authHeader();
    const options ={ headers };
    return requestWrapper(()=>axios.get('users-get', options));
}

const getUser = async ( id ) => {
    const headers= authHeader();
    const options ={ headers };
    return requestWrapper(()=>axios.get(`user/${id}`, options));
};

/**
 * @param {SignUpModel} model
 */
const createUser = async ( model ) => {
    return requestWrapper(()=>axios.post(`users-create`, model ));
};

const validatePassword = async ( password ) => {
    const headers = authHeader();
    const options = { headers };
    return requestWrapper(()=>axios.get(`validate-password/${password}`, options));
};

/**
 * @param {LogInModel} model
 */
const logIn = async ( model ) => {
    return requestWrapper(()=>axios.post(`login`, model));
};

const logOut = async ( model ) => {
    const headers = authHeader();
    const options = { headers };
    return requestWrapper(()=>axios.post(`logout`, model, options));
};

/**
 * @param {Number} id 
 * @param {UpdateUserViewModel} model
 */
const updateUser = async ( id, model ) =>{
    const headers = {
        ...authHeader(),
        'Content-Type': 'multipart/form-data'
    };

    const body = new FormData();
    body.append('username'    , model.username    );
    body.append('tag'         , model.tag         );
    body.append('email'       , model.email       );
    body.append('profilePic'  , model.profilePic  );
    body.append('coverPic'    , model.coverPic    );

    const options = { headers };

    return requestWrapper(()=>axios.put(`users-update/${id}`, body, options));
};

/**
 * @param   {Number} id
 */
const deleteUser = async ( id ) =>{
    const headers = authHeader();
    const options = { headers };
    return requestWrapper(()=>axios.delete(`users-delete/${id}`, options));
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
