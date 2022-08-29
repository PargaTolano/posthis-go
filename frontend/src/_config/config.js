import axios, {AxiosError} from 'axios';
import { authenticationService } from '_services';

const getURL = subroute => `${process.env.REACT_APP_API_HOST}/${subroute}`;

axios.interceptors.response.use(
    response=>response,
    error=>{
    if( [401, 403].indexOf(error.code) !== -1 ){
        authenticationService.logout();
        window.location.reload();
    }
    return Promise.reject(error);
});

const instance = axios.create({baseURL: process.env.REACT_APP_API_HOST});

export {
    getURL,
    instance as axios
}
