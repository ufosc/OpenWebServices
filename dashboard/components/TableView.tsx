// @ts-nocheck
import { TrashCan } from '@carbon/icons-react'
import {
  DataTable, Table, TableHead, TableRow,
  TableHeader, TableBody, TableCell,
  TableContainer, TableToolbar, TableBatchAction,
  TableBatchActions, TableToolbarContent, Button,
  TableToolbarSearch, TableSelectAll, TableSelectRow
} from '@carbon/react'

export default function TableView(props: {
  rows: any, headers: any, title: string, description: string
}) {
  return (
    <DataTable rows={props.rows} headers={props.headers} isSortable>
      {({
        rows,
        headers,
        getTableProps,
        getHeaderProps,
        getRowProps,
        getBatchActionProps,
        selectRow,
        getToolbarProps,
        getSelectionProps
      }) => {

        const batchActionProps = {
          ...getBatchActionProps({
            onSelectAll: () => {
              rows.map(row => {
                if (!row.isSelected) {
                  selectRow(row.id)
                }
              })
            }
          })
        }

        return (
          <TableContainer title={props.title} description={props.description}>
            <TableToolbar {...getToolbarProps()}>
              <TableBatchActions {...batchActionProps}>
                <TableBatchAction tabIndex={batchActionProps.shouldShowBatchActions ? 0 : -1} renderIcon={TrashCan}>
                  Delete
                </TableBatchAction>
              </TableBatchActions>
            </TableToolbar>
            <Table {...getTableProps()}>
              <TableHead>
                <TableRow>
                  <TableSelectAll { ...getSelectionProps() } />
                  {headers.map((header, i) => (
                    <TableHeader key={i} {...getHeaderProps({ header })}>
                      {header.header}
                    </TableHeader>
                  ))}
                </TableRow>
              </TableHead>
              <TableBody>
                {rows.map((row, i) => (
                  <TableRow key={i} {...getRowProps({ row })}>
                    <TableSelectRow {...getSelectionProps({ row })} />
                    {row.cells.map((cell) => (
                      <TableCell key={cell.id}>{cell.value}</TableCell>
                    ))}
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        )}}
    </DataTable>
  )
}
