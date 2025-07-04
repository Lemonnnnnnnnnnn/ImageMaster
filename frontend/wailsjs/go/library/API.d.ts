// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {library} from '../models';
import {context} from '../models';

export function DeleteManga(arg1:string):Promise<boolean>;

export function GetAllMangas():Promise<Array<library.Manga>>;

export function GetImageDataUrl(arg1:string):Promise<string>;

export function GetMangaImages(arg1:string):Promise<Array<string>>;

export function InitializeLibraryManager():Promise<void>;

export function LoadAllLibraries():Promise<void>;

export function SelectLibrary():Promise<string>;

export function SetContext(arg1:context.Context):Promise<void>;

export function SetOutputDir():Promise<string>;
