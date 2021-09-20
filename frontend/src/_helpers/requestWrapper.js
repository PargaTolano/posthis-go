import { handleResponse} from '_helpers';

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
 * @param {Promise<Any>} req request to listen from
 * @returns {Promise<ResponseData>} handled data for use in the front end
 */
export async function requestWrapper(req){
    try{
        let res = await req();
        let handled = await handleResponse(res);
        return new ResponseData( handled, null);
    }catch(e){
        return new ResponseData( null, e);
    }
}