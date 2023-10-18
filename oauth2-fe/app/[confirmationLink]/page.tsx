export default function ConfirmationPage({ params }: { params: { confirmationLink: string } } ) {
    return (
        <p>Confirmation link: {params.confirmationLink}</p>
    )
}