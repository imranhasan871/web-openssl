
// this file is generated — do not edit it


declare module "svelte/elements" {
	export interface HTMLAttributes<T> {
		'data-sveltekit-keepfocus'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-noscroll'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-preload-code'?:
			| true
			| ''
			| 'eager'
			| 'viewport'
			| 'hover'
			| 'tap'
			| 'off'
			| undefined
			| null;
		'data-sveltekit-preload-data'?: true | '' | 'hover' | 'tap' | 'off' | undefined | null;
		'data-sveltekit-reload'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-replacestate'?: true | '' | 'off' | undefined | null;
	}
}

export {};


declare module "$app/types" {
	export interface AppTypes {
		RouteId(): "/" | "/dashboard" | "/dashboard/certificates" | "/dashboard/certificates/convert" | "/dashboard/certificates/csr" | "/dashboard/certificates/generate" | "/dashboard/certificates/keys" | "/dashboard/certificates/parse" | "/dashboard/certificates/verify" | "/dashboard/encryption" | "/dashboard/encryption/asymmetric" | "/dashboard/encryption/hash" | "/dashboard/encryption/symmetric" | "/dashboard/operations" | "/login" | "/register";
		RouteParams(): {
			
		};
		LayoutParams(): {
			"/": Record<string, never>;
			"/dashboard": Record<string, never>;
			"/dashboard/certificates": Record<string, never>;
			"/dashboard/certificates/convert": Record<string, never>;
			"/dashboard/certificates/csr": Record<string, never>;
			"/dashboard/certificates/generate": Record<string, never>;
			"/dashboard/certificates/keys": Record<string, never>;
			"/dashboard/certificates/parse": Record<string, never>;
			"/dashboard/certificates/verify": Record<string, never>;
			"/dashboard/encryption": Record<string, never>;
			"/dashboard/encryption/asymmetric": Record<string, never>;
			"/dashboard/encryption/hash": Record<string, never>;
			"/dashboard/encryption/symmetric": Record<string, never>;
			"/dashboard/operations": Record<string, never>;
			"/login": Record<string, never>;
			"/register": Record<string, never>
		};
		Pathname(): "/" | "/dashboard" | "/dashboard/" | "/dashboard/certificates" | "/dashboard/certificates/" | "/dashboard/certificates/convert" | "/dashboard/certificates/convert/" | "/dashboard/certificates/csr" | "/dashboard/certificates/csr/" | "/dashboard/certificates/generate" | "/dashboard/certificates/generate/" | "/dashboard/certificates/keys" | "/dashboard/certificates/keys/" | "/dashboard/certificates/parse" | "/dashboard/certificates/parse/" | "/dashboard/certificates/verify" | "/dashboard/certificates/verify/" | "/dashboard/encryption" | "/dashboard/encryption/" | "/dashboard/encryption/asymmetric" | "/dashboard/encryption/asymmetric/" | "/dashboard/encryption/hash" | "/dashboard/encryption/hash/" | "/dashboard/encryption/symmetric" | "/dashboard/encryption/symmetric/" | "/dashboard/operations" | "/dashboard/operations/" | "/login" | "/login/" | "/register" | "/register/";
		ResolvedPathname(): `${"" | `/${string}`}${ReturnType<AppTypes['Pathname']>}`;
		Asset(): string & {};
	}
}