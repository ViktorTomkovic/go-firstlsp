local client = vim.lsp.start_client {
	name = "firstLSP",
	cmd = { "~/go/src/github.com/ViktorTomkovic/go-firstlsp/go-firstlsp" },
}

if not client then
	vim.notify('LSP client "firstLSP" is not setup properly')
	return
end

vim.api.nvim_create_autocmd("FileType", {
	pattern = "markdown",
	callback = function()
		vim.lsp.buf_attach_client(0, client)
	end,
})
