import { useContext } from "react";
import { useEditBill } from "../../hooks/admin/useEditBill";
import { useHasPermission } from "../../hooks/admin/useHasPermission";
import { Bill, BillContext } from "../../model/Bill";
import { BillDetailsForm } from "./BillDetailsForm";

export function EditBillForm() {
  const bill = useContext(BillContext);
  const canWrite = useHasPermission("billing", "write");
  const mutation = useEditBill(bill.id);
  return (
    <BillDetailsForm
      initialValues={bill}
      isDisabled={!canWrite}
      onSubmit={async (newBill) => {
        await mutation.mutateAsync(newBill);
      }}
    >
      Save Changes
    </BillDetailsForm>
  );
}