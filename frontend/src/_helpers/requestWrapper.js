import axios, { AxiosResponse } from 'axios';

export class ResponseData{

    data;
    err;

    constructor( data, err ){
        this.data = data;
        this.err = err;
    }
}

/**
 * 
 * @param {Promise<AxiosResponse>} req request to listen from
 * @returns {Promise<ResponseData>} handled data for use in the front end
 */
export async function requestWrapper(req){
    try{
        const res=await req();
        return new ResponseData( res.data, null);
    } catch(e){
        return new ResponseData( null, e);
    }
}